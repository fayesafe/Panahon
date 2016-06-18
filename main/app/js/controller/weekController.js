angular
  .module('weatherApp')
  .controller('WeekController',
    ['$location', '$scope', 'DataService', 'DatetimeService',
    function($location, $scope, DataService, DatetimeService) {

      $scope.default = true;
      $scope.days = [{},{},{},{},{},{}]
      $scope.open = function(day) {
        $location.url('/day?ts=' + day.ts);
      }

      var ts = Date.now();
      for (var i=0; i<6; i++) {
        (function(i) {
          DataService.getAggregatedDataOfDay(ts).then(function(day) {
            $scope.days[i] = day;
          });
        })(i);
        ts = DatetimeService.getLastDayTimestamp(ts);
      }
}]);
