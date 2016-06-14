angular
  .module('weatherApp')
  .directive('chart', ['CHART_OPTIONS', function(CHART_OPTIONS) {
    return {
      scope: {},
      transclude: true,
      templateUrl: 'templates/chart.html',
      controller: 'ChartController',
      link: function link($scope, element, attrs, controller) {

        if (attrs.flow) {
          $scope.flow = 1*attrs.flow;
        }

        if (attrs.zoom) {
          $scope.zoom = attrs.zoom.split(',');
        }

        $scope.chartId = attrs.chartId;
        $scope.title   = attrs.title;
        $scope.style   = attrs.style;

        // set generic id of chart in C3.js options
        $scope.layout = angular.copy(CHART_OPTIONS[attrs.layout]);
        $scope.layout.bindto = '#' + $scope.chartId;

        // set generic id of chart before initializing with C3.js
        $(element).find('div.chart').attr('id', $scope.chartId);

        $scope.generateChart();
      }
    }
  }]);
