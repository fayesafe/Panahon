angular
  .module('weatherApp')
  .controller('RealtimeController',
    ['$interval', '$scope', 'DataService', 'EVENTS',
    function($interval, $scope, DataService, EVENTS) {

      $scope.updateData = function() {
        DataService.getLastData(10).then(function(rows) {
          setTimeout(function(){
            $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
          }, 100);
        });
      };

      $scope.measure = function() {
        $scope.disableMeasurement = true;
        DataService.measure().then(function() {
          $scope.updateData();
          $scope.disableMeasurement = false;
          $scope.$apply();
        });
      };

      $scope.updateData();
  }]);
