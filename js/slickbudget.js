var filters = angular.module('slickbudget.filters', []);
var services = angular.module('slickbudget.services', ['slickbudget.filters']);
var directives = angular.module('slickbudget.directives', ['slickbudget.services', 'slickbudget.filters']);
var controllers = angular.module('slickbudget.controllers', ['slickbudget.services', 'slickbudget.directives', 'slickbudget.filters']);

var slickbudget = angular.module('slickbudget', ['ngAnimate', 'ngRoute', 'ngTouch', 'slickbudget.filters', 
        'slickbudget.services', 'slickbudget.directives', 'slickbudget.controllers']);

slickbudget.config(['$compileProvider', '$controllerProvider','$filterProvider', '$locationProvider', 
        '$logProvider', '$provide', '$routeProvider',
        function($compileProvider, $controllerProvider, $filterProvider, $locationProvider, 
            $logProvider, $provide, $routeProvider) {

    // Available for scripts loaded after the application has been bootstrapped
    slickbudget.compileProvider = $compileProvider;
    slickbudget.controllerProvider = $controllerProvider;
    slickbudget.filterProvider = $filterProvider;
    slickbudget.provide = $provide;
    slickbudget.routeProvider = $routeProvider;

    $logProvider.debugEnabled(false);
    $locationProvider.html5Mode(false);

    $routeProvider.when('/', {
        templateUrl: 'views/home.html',
        controller: 'HomeController as home',
    });

    $routeProvider.when('/track', {
        templateUrl: 'views/track.html',
        controller: 'TrackingController',
    });


    $routeProvider.when('/reports', {
        templateUrl: 'views/reports.html',
        controller: 'ReportsController',
    });

    $routeProvider.when('/preferences', {
        templateUrl: 'views/prefs.html',
        controller: 'PrefsController',
    });

    $routeProvider.otherwise({redirectTo: '/'});

}]);

