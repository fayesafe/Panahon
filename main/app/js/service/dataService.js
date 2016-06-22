angular
  .module('weatherApp')
  .factory('DataService',
    ['$http', '$q', 'DatetimeService', function($http, $q, DatetimeService) {

    return {
      getData: function(series, localStartDate, localEndDate) {
        var deferred = $q.defer();
        var startDate = DatetimeService.toUtcDate(localStartDate);
        var url = '/api/' + series + '/' + DatetimeService.toDateString(startDate);

        if (localEndDate) {
          var endDate = DatetimeService.toUtcDate(localEndDate)
          url += '/' + DatetimeService.toDateString(endDate);
        }

        $.get(url, function(data) {
          var rows = data.Series[0].values;
          rows.forEach(function(row) {
            row[0] = DatetimeService.toLocalTimestamp(row[0]);
            row[5] *= 100; /*rain*/
          });
          rows.splice(0, 0, data.Series[0].columns);
          deferred.resolve(rows);
        });

        return deferred.promise;
      },
      getLastData: function(count) {
        var deferred = $q.defer();
        $.get('/api/sensors/' + count)
          .done(function(response) {
            var rows = response.Series[0].values;
            rows.forEach(function(row) {
              row[3] *= 100; /*rain*/
            });
            rows.splice(0, 0, response.Series[0].columns);
            deferred.resolve(rows);
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
