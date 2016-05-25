angular
  .module('weatherApp')
  .controller('HumidityCtrl', ['$scope', function($scope) {

    $scope.options = {
      data: [
        { temp: 16 }, { temp: 15 }, {temp: 17 } , { temp: 17} , { temp: 5} , { temp: 22} , { temp: 21} , { temp: 7} , { temp: 23} , { temp: 24} , { temp: 24} , { temp: 10} , { temp: 20} , { temp: 6} , { temp: 21} , { temp: 16} , { temp: 5} , { temp: 17} , { temp: 23} , { temp: 15} , { temp: 13 }
      ],
      dimensions: {
        temp: {
          axis: 'y',
            type: 'bar',
            postfix: '%',
            name: 'Humidity'
          }
      }
    };

  }]);
