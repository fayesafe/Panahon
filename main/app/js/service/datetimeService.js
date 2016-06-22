angular
  .module('weatherApp')
  .factory('DatetimeService', function() {

    var localOffset = new Date().getTimezoneOffset()*60*1000;

    return {
      toUtcDate: function(date) {
        return new Date(date.getTime() - localOffset);
      },
      toLocalTimestamp: function(ts) {
        return ts + localOffset;
      },
      toDateString: function(date) {
        return date.toISOString().slice(0, 10);
      },
      toTimeString: function(date) {
        var h = date.getHours();
        var m = date.getMinutes();
        if (h < 10) h = '0' + h;
        if (m < 10) m = '0' + m;
        return h + ':' + m;
      },
      toDayString: function(date) {
        return date.getDate() + '. Tag';
      },
      getNextDay: function(date) {
        var d = new Date(date);
        d.setDate(d.getDate() + 1);
        return d;
      },
      getNextMonth: function(date) {
        var d = new Date(date);
        d.setMonth(d.getMonth() + 1);
        return d;
      },
      getNextYear: function(date) {
        var d = new Date(date);
        d.setFullYear(d.getFullYear() + 1);
        return d;
      },
      getStartOfMonth: function(date) {
        date.setDate(1);
        return date;
      },
      getStartOfYear: function(date) {
        date.setMonth(0);
        return date;
      },
      timestampToUTC: function(ts) {
        return ts - localOffset;
      },
      getStartTimestampOfDay: function(ts) {
        var a = new Date(ts);
        a.setHours(0);
        a.setMinutes(0);
        a.setSeconds(0);
        a.setMilliseconds(0);
        return a.getTime() - localOffset;
      },
      getStartTimestampOfMonth: function(ts) {
        var a = new Date(this.getStartTimestampOfDay(ts));
        a.setDate(1);
        return a.getTime();
      },
      getStartTimestampOfYear: function(ts) {
        var a = new Date(this.getStartTimestampOfMonth(ts));
        a.setMonth(0);
        return a.getTime();
      },
      getLastDayTimestamp: function(ts) {
        var a = new Date(this.getStartTimestampOfDay(ts));
        a.setDate(a.getDate() - 1);
        return a.getTime();
      },
      getNextDayTimestamp: function(ts) {
        var a = new Date(this.getStartTimestampOfDay(ts));
        a.setDate(a.getDate() + 1);
        return a.getTime();
      },
      getNextMonthTimestamp: function(ts) {
        var a = new Date(this.getStartTimestampOfMonth(ts));
        a.setMonth(a.getMonth() + 1);
        return a.getTime();
      },
      getNextYearTimestamp: function(ts) {
        var a = new Date(this.getStartTimestampOfYear(ts));
        console.log(a.toUTCString());
        console.log(a.getFullYear());
        a.setFullYear(a.getFullYear() + 1);
        console.log(a.getYear());
        return a.getTime();
      }
    };
  });
