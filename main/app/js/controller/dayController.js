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

      $scope.datepicker = $('#datepicker').datetimepicker({
        format: 'DD.MM.YYYY',
        defaultDate: new Date(ts)
      }).on('dp.change', function(e) {
        var ts = e.date.valueOf();
        setTimeout(function($scope) {
          $scope.day = DataService.getDataOfDay(ts);
          $scope.$apply();

          $scope.$broadcast(
            EVENTS.DATA_UPDATED,
            $scope.data)
        }, 1000, $scope);
      }).data("DateTimePicker");

      setTimeout(function() {
        DataService.getDataOfDay(ts).then(function(rows) {
          $scope.$broadcast(EVENTS.DATA_UPDATED, rows)
        }, function(response) {
          console.log('Error:',response);
        });


      },600);


}]);
