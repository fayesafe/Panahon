angular
  .module('weatherApp')
  .controller('AllController',
    ['$scope', 'DataService', 'DatetimeService', 'EVENTS',
    function($scope, DataService, DatetimeService, EVENTS) {

      $scope.format = 'DD.MM.YYYY';

      $scope.datepicker = $('#datepicker').datetimepicker({
        format: $scope.format,
        defaultDate: Date.now()
      }).on('dp.change', function(e) {
        var ts = e.date.valueOf();
        $scope.getData(ts);
      }).data("DateTimePicker");

      $scope.$watch('format', function() {
        $scope.datepicker.format($scope.format);

        if ($scope.format == 'DD.MM.YYYY') {
          $scope.getData = DataService.getDataOfDay;
          $scope.formatTimestamp = DatetimeService.formatTime;
        } else if ($scope.format == 'MM.YYYY') {
          $scope.getData = DataService.getDataOfMonth;
          $scope.formatTimestamp = DatetimeService.formatDay;
        } else {
          $scope.getData = DataService.getDataOfYear;
          $scope.formatTimestamp = DatetimeService.formatMonth;
        }

        $scope.loadData($scope.datepicker.date().valueOf());
      });

      $scope.loadData = function(ts) {
        var data = $scope.getData(ts);
        data.forEach(function(elem, index, array) {
          elem.ts = $scope.formatTimestamp(elem.ts);
        });
        $scope.$broadcast(EVENTS.DATA_UPDATED, data);
      };

      setTimeout(function($scope) {
        $scope.loadData($scope.datepicker.date().valueOf());
      }, 1000, $scope);
  }]);
