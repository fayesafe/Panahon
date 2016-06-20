angular
  .module('weatherApp')
  .controller('RealtimeController',
    ['$interval', '$scope', 'DataService', 'EVENTS',
    function($interval, $scope, DataService, EVENTS) {

      $scope.updateData = function() {
        var now = Date.now()
        DataService.getDataBetween($scope.lastUpdate, now).then(function(rows) {
          $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
        });
        $scope.lastUpdate = now;
      };

      $scope.measure = function() {
        $scope.disableMeasurement = true;
        setTimeout(function(){
          $scope.disableMeasurement = false;
          $scope.$apply();
        }, 2000);
        DataService.measure().then(function() {
          $scope.updateData();
        });
      };

      $scope.lastUpdate = Date.now();
      DataService.getLastData(10).then(function(rows) {
        $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
      });
  }]);
