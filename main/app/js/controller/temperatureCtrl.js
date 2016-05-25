angular
  .module('weatherApp')
  .controller('TemperatureCtrl', ['$scope', 'DataService', 'CHARTS', function($scope, DataService, CHARTS) {

    $scope.data = [['time'], ['temp']];

    $scope.options = CHARTS.TEMPERATURE;

    var count = 0;

    $scope.$on('data:updated', function(event, newData) {
      var length = ++count < 10 ? 0 : 1;
      $scope.instance.flow({
        columns: newData, duration: 1000, length: length
      });
    });
  }])
