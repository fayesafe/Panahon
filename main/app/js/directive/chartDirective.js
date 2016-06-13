angular
  .module('weatherApp')
  .directive('chart', function() {
    return {
      scope: {},
      transclude: true,
      templateUrl: 'templates/chart.html',
      controller: 'ChartController',
      link: function link($scope, element, attrs, controller) {

        if (!attrs.flow) {
          attrs.flow = false;
        }
        if (!attrs.flowValuesLimit) {
          attrs.flowValuesLimit = 10;
        }

        $scope.chartId = attrs.chartId;
        $scope.title   = attrs.title;
        $scope.style   = attrs.style;
        $scope.key     = attrs.key;
        $scope.flow    = attrs.flow;
        $scope.flowValuesLimit = 1*attrs.flowValuesLimit;

        var layout = JSON.parse(attrs.layout);
        layout.bindto = '#' + $scope.chartId;

        // set generic id of chart before initializing with C3js
        $(element).find('div.chart').attr('id', $scope.chartId);

        $scope.generateChart(layout);
      }
    }
  });
