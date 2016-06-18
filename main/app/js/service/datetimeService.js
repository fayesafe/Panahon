angular
  .module('weatherApp')
  .factory('DatetimeService',
    ['DATETIME', function(DATETIME) {

    var localOffset = new Date().getTimezoneOffset()*60*1000;

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
        return new Date(ts).getDate();
      },
      formatMonth: function(ts) {
        return DATETIME.MONTHS[new Date(ts).getMonth()];
      },
      getStartTimestampOfDay: function(ts, toUTC) {
        var a = new Date(ts);
        a.setHours(0);
        a.setMinutes(0);
        a.setSeconds(0);
        a.setMilliseconds(0);
        return a.getTime() - (toUTC ? localOffset : 0);
      },
      getStartTimestampOfMonth: function(ts, toUTC) {
        var a = new Date(this.getStartTimestampOfDay(ts, toUTC));
        a.setDate(1);
        return a.getTime();
      },
      getStartTimestampOfYear: function(ts, toUTC) {
        var a = new Date(this.getStartTimestampOfMonth(ts, toUTC));
        a.setMonth(0);
        return a.getTime();
      },
      getLastDayTimestamp: function(ts, toUTC) {
        var a = new Date(this.getStartTimestampOfDay(ts, toUTC));
        a.setDate(a.getDate() - 1);
        return a.getTime() - (toUTC ? localOffset : 0);
      },
      getNextDayTimestamp: function(ts, toUTC) {
        var a = new Date(this.getStartTimestampOfDay(ts, toUTC));
        a.setDate(a.getDate() + 1);
        return a.getTime() - (toUTC ? localOffset : 0);
      },
      getNextMonthTimestamp: function(ts, toUTC) {
        var a = new Date(this.getStartTimestampOfMonth(ts, toUTC));
        a.setMonth(a.getMonth() + 1);
        return a.getTime();
      },
      getNextYearTimestamp: function(ts, toUTC) {
        var a = new Date(this.getStartTimestampOfYear(ts, toUTC));
        console.log(a.toUTCString());
        console.log(a.getFullYear());
        a.setFullYear(a.getFullYear() + 1);
        console.log(a.getYear());
        return a.getTime();
      }
    };
  }]);
