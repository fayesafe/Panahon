angular
  .module('weatherApp')
  .controller('WeekController',
    ['$location', '$scope', 'DataService',
    function($location, $scope, DataService) {

      $scope.default = true;
      $scope.days = DataService.getLastDaysData(6);
      $scope.open = function(day) {
        $location.url('/day?ts=' + day.ts);
      }

}]);
