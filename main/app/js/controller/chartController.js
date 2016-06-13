angular
  .module('weatherApp')
  .controller('ChartController',
    ['$interval', '$scope', 'CHART_OPTIONS', 'EVENTS',
    function($interval, $scope, CHART_OPTIONS, EVENTS) {

      $scope.generateChart = function(layout) {
        $scope.chart = c3.generate(layout);
      };

      $scope.loadData = function(data) {
        $scope.chart.load({
          json: data,
          keys: {
            x: 'time',
            value: [$scope.key]
          }
        });
      };

      $scope.flowData = function(data) {
        var flowLength = 0;
        var existingData = $scope.chart.data.shown();

        // already 10 measurements in chart? If so, then flow existing data out
        if (existingData.length > 0 &&
            existingData[0].values.length > $scope.flowValuesLimit) {
          flowLength = data.length;
        }

        $scope.chart.flow({
          json: data,
          keys: {
            value: ['time', $scope.key]
          },
          length: flowLength
        });
      };

      $scope.$on(EVENTS.DATA_UPDATED, function(event, data) {
        if ($scope.flow) {
          $scope.flowData(data);
        } else {
          $scope.loadData(data);
        }
      });
  }]);
