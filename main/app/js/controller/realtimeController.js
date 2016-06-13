angular
  .module('weatherApp')
  .controller('RealtimeController',
    ['$interval', '$rootScope', '$scope', 'DataService', 'CHART_OPTIONS', 'EVENTS',
    function($interval, $rootScope, $scope, DataService, CHART_OPTIONS, EVENTS) {

      var updateData = function() {
        var data = [{
          time: Date.now(),
          temperature: Math.floor((Math.random() * 25) + -10),
          humidity: Math.floor((Math.random() * 100) + 0),
          pressure: Math.floor((Math.random() * 1200) + 800)
        }];
        $rootScope.$broadcast(EVENTS.DATA_UPDATED, data);
      }

      $scope.timer = $interval(updateData, 2000);

      $scope.$on("$destroy",function(){
        $interval.cancel($scope.timer);
      });

      // start simlation of data updates if site is active and stop if inactive
      $(window).focus(function(){
        $scope.timer = $interval(updateData, 2000);
      }).blur(function(){
        $interval.cancel($scope.timer);
      });
  }]);
