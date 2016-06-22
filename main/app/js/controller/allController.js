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
        $scope.loadData(ts);
      }).data("DateTimePicker");

      $scope.$watch('format', function() {
        $scope.datepicker.format($scope.format);

        if ($scope.format == 'DD.MM.YYYY') {
          $scope.series = 'hours';
          $scope.getStartDate = function(date) { return date; };
          $scope.getEndDate = DatetimeService.getNextDay;
          $scope.formatDate = DatetimeService.toTimeString;
        } else if ($scope.format == 'MM.YYYY') {
          $scope.series = 'days';
          $scope.getStartDate = DatetimeService.getStartOfMonth;
          $scope.getEndDate = DatetimeService.getNextMonth;
          $scope.formatDate = DatetimeService.toDayString;
        } else {
          $scope.series = 'weeks';
          $scope.getStartDate = DatetimeService.getStartOfYear;
          $scope.getEndDate =  DatetimeService.getNextYear;
          $scope.formatDate = function(date, index) {
            return index + '. Woche';
          };
        }

        $scope.loadData(new Date($scope.datepicker.date().valueOf()));
      });

      $scope.loadData = function(startDate) {
        var startDate = $scope.getStartDate(new Date(startDate));
        var endDate = $scope.getEndDate(startDate);

        console.log(startDate, endDate);

        DataService.getData($scope.series, startDate, endDate).then(function(rows) {
          console.log(rows);
          rows.forEach(function(elem, index, array) {
            if (index != 0) {
              elem[0] = $scope.formatDate(new Date(elem[0]), index);
            }
          });
          $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
        }, function(response) {
          console.log('Error:',response);
        });
      };
  }]);
