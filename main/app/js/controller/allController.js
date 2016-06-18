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
        console.log('change',new Date(ts).toString());
        $scope.loadData(ts);
      }).data("DateTimePicker");

      $scope.$watch('format', function() {
        $scope.datepicker.format($scope.format);

        if ($scope.format == 'DD.MM.YYYY') {
          $scope.getData = DataService.getDataOfDay;
          $scope.formatTimestamp = DatetimeService.formatTime;
          $scope.datepicker.date()
        } else if ($scope.format == 'MM.YYYY') {
          $scope.getData = DataService.getDataOfMonth;
          $scope.formatTimestamp = DatetimeService.formatDay;
        } else {
          $scope.getData = DataService.getDataOfYear;
          $scope.formatTimestamp = function(ts, index) { return index+'. Woche'};
        }

        $scope.loadData($scope.datepicker.date().valueOf());
      });

      $scope.loadData = function(ts) {
        $scope.getData(ts).then(function(rows) {
          rows.forEach(function(elem, index, array) {
            if (index != 0) {
              elem[0] = $scope.formatTimestamp(elem[0], index);
            }
          });
          $scope.$broadcast(EVENTS.DATA_UPDATED, rows);
        }, function(response) {
          console.log('Error:',response);
        });
      };

      /*setTimeout(function($scope) {
        $scope.loadData($scope.datepicker.date().valueOf());
      }, 1000, $scope);*/
  }]);
