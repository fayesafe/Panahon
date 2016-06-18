angular
  .module('weatherApp')
  .controller('ChartController',
    ['$interval', '$scope', 'CHART_OPTIONS', 'EVENTS',
    function($interval, $scope, CHART_OPTIONS, EVENTS) {

      $scope.generateChart = function() {
        $scope.chart = c3.generate($scope.layout);
      };

      $scope.loadData = function(rows) {
        //$scope.chart.unload({ids: $scope.layout.data.keys});
        $scope.chart.load({
          rows: rows
        });

        if ($scope.zoom) {
          $scope.chart.zoom($scope.zoom);
        }
      };

      $scope.flowData = function(rows) {
        var flowLength = 0;
        var existingData = $scope.chart.data.shown();

        // already 10 measurements in chart? If so, then flow existing data out
        if (existingData.length > 0 &&
            existingData[0].values.length > $scope.flow) {
          flowLength = rows.length;
        }

        $scope.chart.flow({
          rows: rows,
          length: flowLength
        });
      };

      $scope.$on(EVENTS.DATA_UPDATED, function(event, rows) {
        if ($scope.flow) {
          $scope.flowData(rows);
        } else {
          $scope.loadData(rows);
        }
      });
  }]);
