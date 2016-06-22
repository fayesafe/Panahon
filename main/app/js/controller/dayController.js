angular
  .module('weatherApp')
  .controller('DayController',
    ['$routeParams', '$scope', 'DataService', 'EVENTS',
    function($routeParams, $scope, DataService, EVENTS) {

      var date = new Date();
      if ($routeParams.ts) {
        date = new Date(1*$routeParams.ts);
      }

      $scope.loadDay = function(date) {
        DataService.getData('hours', date).then(function(rows) {
          setTimeout(function() {
            console.log(rows);
            $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
          }, 100);
        }, function(response) {
          console.log('Error:',response);
        });

        DataService.getData('days', date).then(function(days) {
          $scope.day = days[1];
          console.log(days);
        }, function(response) {
          console.log('Error:',response);
        });
      };

      $scope.datepicker = $('#datepicker').datetimepicker({
        format: 'DD.MM.YYYY',
        defaultDate: date
      }).on('dp.change', function(e) {
        var ts = e.date.valueOf();
        $scope.loadDay(new Date(ts));
      }).data("DateTimePicker");

      $scope.loadDay(date);
}]);
