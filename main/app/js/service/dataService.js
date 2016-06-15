angular
  .module('weatherApp')
  .factory('DataService',
    ['$http', '$q', 'DatetimeService', function($http, $q, DatetimeService) {

    return {
      getDataOfYear: function(ts) {
        var data = [];

        for (var i=0; i<12; i++) {
          data.push({
            ts: ts + 100000*i,
            temperature: Math.floor((Math.random() * 25) + -10),
            humidity: Math.floor((Math.random() * 100) + 0),
            pressure: Math.floor((Math.random() * 1200) + 800),
            rain: Math.floor((Math.random() * 100) + 0),
            sunHours: Math.floor((Math.random() * 12) + 5)
          });
        }

        return data;
      },
      getDataOfMonth: function(ts) {
        var data = [];

        for (var i=0; i<30; i++) {
          data.push({
            ts: ts + 100000*i,
            temperature: Math.floor((Math.random() * 25) + -10),
            humidity: Math.floor((Math.random() * 100) + 0),
            pressure: Math.floor((Math.random() * 1200) + 800),
            rain: Math.floor((Math.random() * 100) + 0),
            sunHours: Math.floor((Math.random() * 12) + 5)
          });
        }

        return data;
      },
      getDataOfDay: function(ts) {

        var deferred = $q.defer();
        var tsStart = DatetimeService.getStartTimestampOfDay(ts);
        var tsEnd = DatetimeService.getNextDayTimestamp(ts);

        $http.get('/api/range/' + tsStart + '/' + tsEnd).then(function(response) {
          var rows = response.data.Series[0].values;
          rows.splice(0, 0, response.data.Series[0].columns);
          deferred.resolve(rows);
        }, function(response) {
          deferred.reject(response);
        });

        return deferred.promise;
      },
      getAggregatedDataOfDay: function(ts) {
        return {
          ts: ts,
          temperatureMax: Math.floor((Math.random() * 25) + -10),
          temperatureMin: Math.floor((Math.random() * 25) + -10),
          humidity: Math.floor((Math.random() * 100) + 0),
          pressure: Math.floor((Math.random() * 1200) + 800),
          rain: Math.floor((Math.random() * 100) + 0),
          sunHours: Math.floor((Math.random() * 12) + 5)
        };
      },
      getAggregatedDataOfLastDays: function(dayCount) {
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
