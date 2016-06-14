angular
  .module('weatherApp')
  .controller('DayController',
    ['$routeParams', '$scope', 'DataService', 'EVENTS',
    function($routeParams, $scope, DataService, EVENTS) {

      var ts = Date.now();
      if ($routeParams.ts) {
        ts = 1*$routeParams.ts;
      }

      $scope.day = {};
      $scope.day.ts = ts;

      setTimeout(function($scope) {
        $scope.day = DataService.getDayData(ts);
        $scope.$apply();

        $scope.$broadcast(
          EVENTS.DATA_UPDATED,
          $scope.day.data)
      }, 1000, $scope);

      $scope.datepicker = $('#datepicker').datetimepicker({
        format: 'DD.MM.YYYY',
        defaultDate: new Date(ts)
      });
      $scope.datepicker.on('changeDate', function(e) {
        console.log(e);
      });
}]);
