angular
  .module('weatherApp')
  .factory('DataService',
    ['$http', '$q', 'DatetimeService', function($http, $q, DatetimeService) {

    var getData = function(tsStart, tsEnd, interval) {
      var deferred = $q.defer();

      $.when(
        $.get('/api/av/temperature/'+interval+'/'+tsStart+'/'+tsEnd),
        $.get('/api/av/pressure/'+interval+'/'+tsStart+'/'+tsEnd),
        $.get('/api/av/humidity/'+interval+'/'+tsStart+'/'+tsEnd),
        $.get('/api/av/rain/'+interval+'/'+tsStart+'/'+tsEnd),
        $.get('/api/av/sun/'+interval+'/'+tsStart+'/'+tsEnd)
      ).done(function(avTemp, avPressure, avHumidity, avRain, avSun) {
        var rows = [['time','temperature','pressure','humidity','rain','sun']];
        console.log('avTemp',avTemp);
        for (var i=1; i<avTemp[0].Series[0].values.length; i++) {
          rows[i] = [
            avTemp[0].Series[0].values[i][0],
            avTemp[0].Series[0].values[i][1],
            avPressure[0].Series[0].values[i][1],
            avHumidity[0].Series[0].values[i][1],
            avRain[0].Series[0].values[i][1] * 100,
            avSun[0].Series[0].values[i][1]
          ];
        };
        deferred.resolve(rows);
      });

      return deferred.promise;
    };

    return {
      getDataOfYear: function(ts) {
        var tsStart = DatetimeService.getStartTimestampOfYear(ts);
        var tsEnd = DatetimeService.getNextYearTimestamp(ts);
        console.log(new Date(tsStart).toUTCString(),new Date(tsEnd).toUTCString());
        return getData(tsStart, tsEnd, '1w');
      },
      getDataOfMonth: function(ts) {
        var tsStart = DatetimeService.getLastDayTimestamp(
          DatetimeService.getStartTimestampOfMonth(ts));
        var tsEnd = DatetimeService.getNextMonthTimestamp(ts);
        console.log(new Date(tsStart).toUTCString(),new Date(tsEnd).toUTCString());
        return getData(tsStart, tsEnd, '1d');
      },
      getDataOfDay: function(ts) {
        var tsStart = DatetimeService.getStartTimestampOfDay(ts) - 3600000;
        var tsEnd = DatetimeService.getNextDayTimestamp(ts);
        console.log(new Date(tsStart).toUTCString(),new Date(tsEnd).toUTCString());
        return getData(tsStart, tsEnd, '1h');
      },
      getAggregatedDataOfDay: function(ts) {
        var deferred = $q.defer();
        var tsStart = DatetimeService.getStartTimestampOfDay(ts);
        var tsEnd = DatetimeService.getNextDayTimestamp(ts);

        $.when(
          $.get('/api/max/temperature/1d/'+tsStart+'/'+tsEnd),
          $.get('/api/min/temperature/1d/'+tsStart+'/'+tsEnd),
          $.get('/api/av/pressure/1d/'+tsStart+'/'+tsEnd),
          $.get('/api/av/humidity/1d/'+tsStart+'/'+tsEnd),
          $.get('/api/av/rain/1d/'+tsStart+'/'+tsEnd),
          $.get('/api/av/sun/1h/'+tsStart+'/'+tsEnd)
        ).done(function(maxTemp, minTemp, avPressure, avHumidity, avRain, avSun) {
          var sunHours = 0;
          avSun[0].Series[0].values.forEach(function(elem) {
            if (elem[1] < 300) sunHours += 1;
          });

          deferred.resolve({
            ts: ts,
            temperatureMax: maxTemp[0].Series[0].values[0][1],
            temperatureMin: minTemp[0].Series[0].values[0][1],
            pressure: avPressure[0].Series[0].values[0][1],
            humidity: avHumidity[0].Series[0].values[0][1],
            rain: avRain[0].Series[0].values[0][1] * 100,
            sunHours: sunHours
          });
        }).fail(function(response) {
          deferred.reject(response);
        });

        return deferred.promise;
      },
      getLastData: function(count) {
        var deferred = $q.defer();
        $.get('/api/last/'+count)
          .done(function(response) {
            var data = response.Series[0].values;
            data.splice(0, 0, response.Series[0].columns);
            deferred.resolve(data);
          }).fail(function(response) {
            deferred.reject(response);
          });
        return deferred.promise;
      },
      getDataBetween: function(tsStart, tsEnd) {
        var deferred = $q.defer();
        $.get('/api/range/'+DatetimeService.timestampToUTC(tsStart)+'/'+DatetimeService.timestampToUTC(tsEnd))
          .done(function(response) {
            if (response.Series) {
              deferred.resolve(response.Series[0].values);
            } else {
              deferred.reject(response);
            }
          }).fail(function(response) {
            deferred.reject(response);
          });
        return deferred.promise;
      },
      measure: function() {
        var deferred = $q.defer();
        $.get('/api/measure')
          .done(function(response) {
            deferred.resolve(response);
          }).fail(function(response) {
            deferred.rejectresponse();
          });
        return deferred.promise;
      }
    };
  }]);
