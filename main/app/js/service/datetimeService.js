angular
  .module('weatherApp')
  .factory('DatetimeService',
    ['DATETIME', function(DATETIME) {

    return {
      formatTime: function(ts) {
        var a = new Date(ts);
        var h = a.getHours();
        var m = a.getMinutes();
        if (h<10) h = '0' + h;
        if (m<10) m = '0' + m;
        return h + ':' + m;
      },
      formatDay: function(ts) {
        return new Date(ts).getDay();
      },
      formatMonth: function(ts) {
        return DATETIME.MONTHS[new Date(ts).getMonth()];
      },
      getStartTimestampOfDay: function(ts) {
        var a = new Date(ts);
        a.setHours(0);
        a.setMinutes(0);
        a.setSeconds(0);
        a.setMilliseconds(0);
        return a.getTime();
      },
      getNextDayTimestamp: function(ts) {
        var a = new Date(ts);
        a.setDate(a.getDate() + 1)
        return a.getTime();
      }
    };
  }]);
