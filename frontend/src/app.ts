/// <reference path="../typings/tsd.d.ts" />

module mailerLogApp {
	'use strict';

	require('angular');
	require('angular-route');
    require('angular-resource');
	require('./scss/index.scss');
    
    var mailerLogControllers = angular.module('mailerLogControllers', []);
    
    mailerLogControllers.controller('mailerTableController', ['$scope', '$location', 'Mail', 
    ($scope, $location, Mail) => {
        $scope.perPage = 10;
        $scope.curPage = 1;
        
        $scope.getPage = () => {
            Mail.query({
                page: $scope.curPage,
                per_page: $scope.perPage
            }, 
            (mails, responseHeaders) => {
                $scope.mails = mails;
                $scope.totalPages = Math.ceil(responseHeaders('X-Total-Count') / $scope.perPage);
            });  
        };
        
        $scope.getPage();
        
        $scope.nextPage = () => {
            $scope.curPage += 1;
            $scope.getPage();
        };
        
        $scope.previousPage = () => {
            $scope.curPage -= 1;
            $scope.getPage();
        };

        $scope.itemClick = (mailId: string) => {
            $location.path('/mail/' + mailId);
        };
        
    }]);
 
    mailerLogControllers.controller('mailerDetailController', ['$scope', '$route', 'Mail', 
    ($scope, $route, Mail) => {
        var params = $route.current.params;
        $scope.mail = Mail.get({mailId: params.mailId});  
    }]);   
    
    var mailerLogServices = angular.module('mailerLogServices', ['ngResource']);
    
    mailerLogServices.factory('Mail', ['$resource', 
    ($resource) => {
        return $resource('mails/:mailId', {}, {
            query: {method: 'GET', params: {mailId: '@id'}, isArray: true}
        });
    }]);


	var app = angular.module('mailer-log-app', ['ngRoute', 'mailerLogControllers', 'mailerLogServices']);

	app.config(['$routeProvider',
    ($routeProvider: ng.route.IRouteProvider) => {
        $routeProvider
            .when('/', {
                template: require('./templates/mailer-table.html'),
                controller: 'mailerTableController',
                name: 'mail-list'
            })
            .when('/mail/:mailId', {
                template: require('./templates/mailer-detail.html'),
                controller: 'mailerDetailController',
                name: 'mail-detail'
            })
            .otherwise({
                redirectTo: '/'
            });
    }]);
    
    console.log('end');
}