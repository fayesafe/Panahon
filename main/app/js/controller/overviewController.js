angular
  .module('weatherApp')
  .controller('OverviewController',
    ['$interval', '$rootScope', '$scope', 'DataService', 'EVENTS',
    function($interval, $rootScope, $scope, DataService, EVENTS) {

      var data = [];
      for (var i=0; i<=20; i++) {
        data.push({
          time: Date.now() + 10000*i,
          temperature: Math.floor((Math.random() * 25) + -10),
          humidity: Math.floor((Math.random() * 100) + 0),
          pressure: Math.floor((Math.random() * 1200) + 800)
        });
      }

      setTimeout(function() {
        $rootScope.$broadcast(EVENTS.DATA_UPDATED, data)
      }, 2000);
}]);
