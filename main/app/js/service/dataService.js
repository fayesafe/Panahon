angular
  .module('weatherApp')
  .factory('DataService',
    ['$rootScope', '$interval', 'EVENTS',
    function($rootScope, $interval, EVENTS) {

    var lastUpdate = new Date();
    var stopInterval;

    var getSensorData = function(since) {
      // var data = call api with last timestamp
      var data = [ ['time'], ['temp'] ];
      data[0].push(Date.now());
      data[1].push(Math.floor((Math.random() * 30) - 5));

      $rootScope.$broadcast(EVENTS.DATA_UPDATED, data);
    };

    // start simlation of data updates if site is active
    $(window).focus(function(){
      stopInterval = $interval(getSensorData, 2000, lastUpdate);
    });

    // stop simlation of data updates if site is inactive
    $(window).blur(function(){
      lastUpdate = new Date();
      $interval.cancel(stopInterval);
    });

    // kick off the data update simlation
    stopInterval = $interval(getSensorData, 2000, lastUpdate);

    return this;
  }]);
