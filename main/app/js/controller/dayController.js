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
        $scope.day = DataService.getDataOfDay(ts);
        $scope.$apply();

        $scope.$broadcast(
          EVENTS.DATA_UPDATED,
          $scope.day.data)
      }, 1000, $scope);

      $('#datepicker').datetimepicker({
        format: 'DD.MM.YYYY',
        defaultDate: new Date(ts)
      }).on('dp.change', function(e) {
        var ts = e.date.valueOf();
        setTimeout(function($scope) {
          $scope.day = DataService.getDataOfDay(ts);
          $scope.$apply();

          $scope.$broadcast(
            EVENTS.DATA_UPDATED,
            $scope.day.data)
        }, 1000, $scope);
      });
}]);
