angular
  .module('weatherApp')
  .controller('TemperatureCtrl',
    ['$scope', 'DataService', 'CHARTS', 'EVENTS',
    function($scope, DataService, CHARTS, EVENTS) {

    $scope.data = [['time'], ['temp']];
    $scope.options = CHARTS.TEMPERATURE;

    $scope.$on(EVENTS.DATA_UPDATED, function(event, data) {
      var flowLength = 0;
      var existingData = $scope.instance.data();

      // already 10 measurements in chart? If so, then flow existing data out
      if (existingData.length > 0 && existingData[0].values.length > 10) {
        flowLength = data[0].length - 1;
      }

      $scope.instance.flow({
        columns: data, duration: 1000, length: flowLength
      });
    });
  }])
