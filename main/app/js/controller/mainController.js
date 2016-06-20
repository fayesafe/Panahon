angular
  .module('weatherApp')
  .controller('MainController',
    ['$scope', function($scope) {
      setTimeout(function(){
        particlesJS.load('footer', './assets/particlesjs-config.json', function() {
          console.log('callback - particles.js config loaded');
        });
      }, 2000);
    }
]);
