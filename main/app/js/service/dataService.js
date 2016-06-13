angular
  .module('weatherApp')
  .factory('DataService',
    ['$rootScope', '$interval', 'EVENTS',
    function($rootScope, $interval, EVENTS) {

    var lastUpdate = new Date();
    var stopInterval;

    var getDataSince = function(since) {
      // var data = call api with last timestamp

      var data = [{
        time: Date.now(),
        temp: Math.floor((Math.random() * 30) - 5)
      }];

      return data;
    };

    return this;
  }]);
