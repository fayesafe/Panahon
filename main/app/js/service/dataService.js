angular
  .module('weatherApp')
  .factory('DataService', ['$rootScope', '$interval', function($rootScope, $interval) {

    $interval(function() {
      var newData = [
        ['time', Date.now()],
        ['temp', Math.floor((Math.random() * 30) - 5)]
      ];
      $rootScope.$broadcast('data:updated', newData);
      console.log('notified');
    }, 1500);

    return this;
  }]);
