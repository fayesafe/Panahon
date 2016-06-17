angular
  .module('weatherApp')
  .controller('DayController',
    ['$routeParams', '$scope', 'DataService', 'EVENTS',
    function($routeParams, $scope, DataService, EVENTS) {

      var ts = Date.now();
      if ($routeParams.ts) {
        ts = 1*$routeParams.ts;
      }

      $scope.loadDay = function(ts) {
        DataService.getDataOfDay(ts).then(function(rows) {
          setTimeout(function() {
            $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
          }, 100);
        }, function(response) {
          console.log('Error:',response);
        });

        DataService.getAggregatedDataOfDay(ts).then(function(day) {
          $scope.day = day;
        }, function(response) {
          console.log('Error:',response);
        });
      };

      $scope.datepicker = $('#datepicker').datetimepicker({
        format: 'DD.MM.YYYY',
        defaultDate: new Date(ts)
      }).on('dp.change', function(e) {
        var ts = e.date.valueOf();
        $scope.loadDay(ts);
      }).data("DateTimePicker");

      $scope.loadDay(ts);
}]);
