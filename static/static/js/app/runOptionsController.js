(function () {

    angular.module('baseApp')
        .controller('runOptionsController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', runOptionsController]);

    function runOptionsController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.defObj = {};
      $scope.loadingParams = true;

      self.initModal = function(){
        $('#runOptionsModal-tagsDropdown').dropdown({ allowAdditions: true, });

        $('#runOptionsModal').modal({//setup button callbacks + general parameters
          closable: false,
          autofocus: false,
          dimmerSettings: {
            closable: false,
            opacity: 0,
          },
          onApprove: self.submitPressed,
          onDeny: self.cancelPressed
        });
      }

      self.incrementVersion = function(){ //set the version to the last version plus one.
        console.log($scope.defObj);
        if ($scope.defObj['last-version'] == null || $scope.defObj['last-version'] == undefined || $scope.defObj['last-version'] == ""){
          $scope.defObj['last-version'] = "0.0.1";
        } else {
          $scope.defObj['last-version'] = $scope.defObj['last-version'].replace(/\d+$/, function(n){ return ++n });
        }
        console.log($scope.defObj);
      }

      self.setupAndShow = function(definitionObj){
        $scope.defObj = jQuery.extend({}, definitionObj);//shallow copy
        self.incrementVersion();
        self.initModal();
        setTimeout(function(){self.initFields();}, 200);
        $('#runOptionsModal').modal('show');
      }
      var startListener = $rootScope.$on('runOptionsModal-start', function(event, defObj) {
        self.setupAndShow(defObj);
      });

      $scope.$on('$destroy', function() {
        startListener();
      });


      //called just before modal show to initialize any javascript on field views.
      self.initFields = function(){
        for (var i = 0; i < $scope.defObj.params.length; i++)
        {
          if ($scope.defObj.params[i].type == "branchselect"){
            $('#runopt-field-' + i).dropdown({
              apiSettings: {
                url: '/api/lookup/buildparam?name=' + $scope.defObj.name + '&param=' + i
              }
            });
          }
        }
        $scope.loadingParams = false;
      }



      //returns the namespace of user-selected build parameters. Defaults are set for untouched fields.
      self.getBuildParamsValues = function(){
        var parameters = {};
        for (var i = 0; i < $scope.defObj.params.length; i++)
        {
          if ($scope.defObj.params[i].type == "text"){
            parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].default;
          }
          if ($scope.defObj.params[i].type == "check"){
            parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].default ? "true" : "false";
          }
          if ($scope.defObj.params[i].type == "branchselect"){
            parameters[$scope.defObj.params[i].varname] = $('#runopt-field-' + i).dropdown('get value');
          }
        }
        console.log("buildParameters: ", parameters);
        return parameters;
      }


      //callbacks called when the user presses a button
      self.submitPressed = function(){
        var tagArray = $('#runOptionsModal-tagsDropdown').dropdown('get value').split(",")
        var isPhysDisabled = $("#runOptionsModal-disablephys").prop('checked');
        var version = $scope.defObj['last-version'];
        $('#runOptionsModal').modal('hide');
        $rootScope.$broadcast('runOptionsModal-finished', {version: version, isPhysDisabled: isPhysDisabled, tags: tagArray, defObj: $scope.defObj, params: self.getBuildParamsValues()});
      }
      self.cancelPressed = function(){
        $('#runOptionsModal').modal('hide');
      }
    }

})();
