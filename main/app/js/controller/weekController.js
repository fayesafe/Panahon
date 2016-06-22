angular
  .module('weatherApp')
  .controller('WeekController',
    ['$location', '$scope', 'DataService', 'DatetimeService',
    function($location, $scope, DataService, DatetimeService) {

      $scope.default = true;
      $scope.days = [{},{},{},{},{},{}]
      $scope.open = function(ts) {
        $location.url('/day?ts=' + ts);
      }

      var today = new Date();
      today.setDate(today.getDate() - 30); // TODO
      var sixDaysAgo = new Date(today);
      sixDaysAgo.setDate(sixDaysAgo.getDate() - 6);

      DataService.getData('days', sixDaysAgo, today).then(function(days) {
        days.splice(0, 1);
        $scope.days = days;
      });
}]);
