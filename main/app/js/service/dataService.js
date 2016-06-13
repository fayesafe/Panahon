angular
  .module('weatherApp')
  .factory('DataService',
    ['$rootScope', '$interval', 'EVENTS',
    function($rootScope, $interval, EVENTS) {

    return {
      getDayData: function(ts) {
        var day = {
          ts: ts,
          temperatureMax: Math.floor((Math.random() * 25) + -10),
          temperatureMin: Math.floor((Math.random() * 25) + -10),
          humidity: Math.floor((Math.random() * 100) + 0),
          pressure: Math.floor((Math.random() * 1200) + 800),
          rain: Math.floor((Math.random() * 100) + 0),
          sunHours: Math.floor((Math.random() * 12) + 5),
          data: []
        };

        for (var i=0; i<24; i++) {
          day.data.push({
            ts: ts + 10000*i,
            temperature: Math.floor((Math.random() * 25) + -10),
            humidity: Math.floor((Math.random() * 100) + 0),
            pressure: Math.floor((Math.random() * 1200) + 800),
            rain: Math.floor((Math.random() * 100) + 0),
            sunHours: Math.floor((Math.random() * 12) + 5)
          });
        }

        return day;
      },
      getLastDaysData: function(dayCount) {
        var days = [];
        for (var i=0; i<dayCount; i++) {
          days.push({
            ts: Date.now() + 10000*i,
            temperatureMax: Math.floor((Math.random() * 25) + -10),
            temperatureMin: Math.floor((Math.random() * 25) + -10),
            humidity: Math.floor((Math.random() * 100) + 0),
            pressure: Math.floor((Math.random() * 1200) + 800),
            rain: Math.floor((Math.random() * 100) + 0),
            sunHours: Math.floor((Math.random() * 12) + 5)
          });
        }
        return days;
      }
    };
  }]);
