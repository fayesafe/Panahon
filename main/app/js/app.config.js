angular.
  module('weatherApp').
  config(['$locationProvider' ,'$routeProvider',
    function config($locationProvider, $routeProvider) {
      $locationProvider.hashPrefix('!');

      $routeProvider
        .when('/week', {
          templateUrl : 'templates/week.html',
          controller: 'WeekController'
        })
        .when('/day', {
          templateUrl : 'templates/day.html',
          controller: 'DayController'
        })
        .when('/realtime', {
          templateUrl : 'templates/realtime.html',
          controller: 'RealtimeController'
        })
        .when('/all', {
          templateUrl : 'templates/all.html',
          controller: 'AllController'
        })
        .otherwise('/week');
    }
  ]);
