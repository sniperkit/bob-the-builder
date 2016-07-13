package builder

import (
  "bobthebuilder/logging"
  "bobthebuilder/util"
  "encoding/json"
  "io/ioutil"
  "path"
  "os"
)

const BASE_FOLDER_NAME = "base"

type BuildDefinition struct {
  Name string `json:"name"`
  Icon string `json:"icon,omitempty"`
  AptPackagesRequired []string `json:"apt-packages-required,omitempty"`
  BaseFolder string `json:"base-folder,omitempty"`
  GitSrc string `json:"git-src,omitempty"`
  Steps []struct {
    Type string `json:"type"`
    Command string `json:"command,omitempty"`
    CanFail bool `json:"can-fail,omitempty"`
    Args []string `json:"args,omitempty"`


    FileName string `json:"filename,omitempty"`
    DestinationFileName string  `json:"filename-destination,omitempty"`
    Directories []string  `json:"directories,omitempty"`
    Files []string  `json:"files,omitempty"`

    Key string `json:"key,omitempty"`
    Value string `json:"value,omitempty"`

    //used for S3 commands
    Bucket string `json:"bucket,omitempty"`
    Region string `json:"region,omitempty"`
    ACL string `json:"ACL,omitempty"`
  } `json:"steps"`

  AbsolutePath string `json:"-"`

  //stateful information
  LastVersion string `json:"last-version"`
  LastRunTime int64 `json:"last-run-time"`




  Params []struct {
    Type string `json:"type"`
    Label string `json:"label"`
    Varname string `json:"varname"`
    Placeholder string `json:"placeholder"`
    Default interface{} `json:"default"`
  } `json:"params"`






  CurrentRun *Run
}






func (d *BuildDefinition)Validate()bool{
  if d.BaseFolder != "" {
    pwd, _ := os.Getwd() //cant have error - would have failed in file/util.go
    if _, err := os.Stat(path.Join(pwd, BASE_FOLDER_NAME, d.BaseFolder)); os.IsNotExist(err) {
      logging.Error("definition-validate", d.Name + " base folder does not exist.")
      return false// base folder does not exist
    }
  }
  return true
}



func (d *BuildDefinition)Flush()error{
  var temp BuildDefinition //make a copy so we can set currentRun to nil
  temp = *d
  temp.CurrentRun  = nil

  data, err := json.MarshalIndent(&temp, "", "  ")
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(temp.AbsolutePath, data, 777)
  return err
}




func (d *BuildDefinition)genRun(tags []string, version string, physDisabled bool)*Run{
  if version == "" || version == "?" {
    version = "0.0.1"
  }

  out := &Run{
    Definition: d,
    GUID: util.RandAlphaKey(32),
    ExecType: "BUILD",
    Version: version,
    Status: STATUS_NOT_YET_RUN,
    Tags: tags,
    PhysDisabled: physDisabled,
  }
  pwd, _ := os.Getwd() //cant have error - would have failed in file/util.go

  //generate phase to clean up the build folder
  delPhase := &CleanPhase{
    DeletePath: path.Join(pwd, BUILD_TEMP_FOLDER_NAME),
  }
  delPhase.init(len(out.Phases))
  out.Phases = append(out.Phases, delPhase)

  if len(d.AptPackagesRequired) > 0 {
    p := &AptGetCheckInstallPhase{ //generate phase to copy in files
      Packages: d.AptPackagesRequired,
    }
    p.init(len(out.Phases))
    out.Phases = append(out.Phases, p)
  }

  if d.GitSrc != "" {//if we are sourcing files from git, that needs to happen first for reasons.
    p := &GitClonePhase{ //generate phase to clone in a git repo
      GitSrcPath: d.GitSrc,
    }
    p.init(len(out.Phases))
    out.Phases = append(out.Phases, p)
  }

  if d.BaseFolder != "" {//next, copy in any static files specified.
    p := &BaseInstallPhase{ //generate phase to copy in files
      BaseAbsPath: path.Join(pwd, BASE_FOLDER_NAME, d.BaseFolder),
    }
    p.init(len(out.Phases))
    out.Phases = append(out.Phases, p)
  }


  //generate a corresponding phase for each step
  for _, step := range d.Steps {
    switch step.Type{
    case "CMD":
      cmd := &CommandPhase{
        Command: step.Command,
        Args: step.Args,
        CanFail: step.CanFail,
      }
      cmd.init(len(out.Phases))
      out.Phases = append(out.Phases, cmd)

    case "EXEC":
      cmd := &ScriptExecPhase{
        ScriptPath: step.Command,
        CanFail: step.CanFail,
      }
      cmd.init(len(out.Phases))
      out.Phases = append(out.Phases, cmd)

    case "S3_UPLOAD":
      cmd := &S3UploadPhase{
        Bucket: step.Bucket,
        Region: step.Region,
        FilenameToUpload: step.FileName,
        DestinationFileName: step.DestinationFileName,
      }
      cmd.init(len(out.Phases), step.ACL)
      out.Phases = append(out.Phases, cmd)

    case "ENV_SET":
      cmd := &SetEnvPhase{
        Key: step.Key,
        Value: step.Value,
      }
      cmd.init(len(out.Phases))
      out.Phases = append(out.Phases, cmd)

    case "TAR_TO_S3":
      cmd := &TarToS3{
        Bucket: step.Bucket,
        Region: step.Region,
        DestinationPath: step.DestinationFileName,
        Directories: step.Directories,
        Files: step.Files,
      }
      cmd.init(len(out.Phases))
      out.Phases = append(out.Phases, cmd)
    }
  }


  return out
}
