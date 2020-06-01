var app = angular.module('printJiraApp', []);
app.controller("mainController", function ($scope, $http) {
  $http.get('./data.json').then(response => {
    $scope.stories = response.data;
  })
});

app.filter('avatar', function () {
  const USER = {
    NICK: {
      name: 'NICK',
      url: 'https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/5be3fb6fa430e70fb00dfecd/c2d642a8-5ef0-4371-b230-00443443b7d5/128?size=48&s=48'
    },
    NAMT: {
      name: 'NAMT',
      url: 'https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/5be3fb6fd9c83f0a32a117fa/44b02a36-e44e-4f8e-bc11-ee15aa133db6/128?size=48&s=48'
    },
    TOAN: {
      name: 'TOAN',
      url: 'https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/5d39064d6508dd0c25fe3940/52a188ef-aeac-4d86-91d8-46eda6bf7228/128?size=48&s=48'
    },
    CONG: {
      name: 'CONG',
      url: 'https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/5d649404e90c310c1722276b/e044cfe4-bba4-4e3b-bb0a-7e49f4d51663/128?size=48&s=48'
    }
  };
  return function (input) {
    switch (input.toUpperCase()) {
      case USER.NICK.name:
        return USER.NICK.url;
      case USER.NAMT.name:
        return USER.NAMT.url;
      case USER.TOAN.name:
        return USER.TOAN.url;
      case USER.CONG.name:
        return USER.CONG.url;
      default:
        return '';
    }
  }
});

app.filter('type', function () {
  const type = {
    STORY: {
      name: 'STORY',
      url: 'https://yuriqa.atlassian.net/secure/viewavatar?size=medium&avatarId=10315&avatarType=issuetype'
    },
    BUG: {
      name: 'BUG',
      url: 'https://yuriqa.atlassian.net/secure/viewavatar?size=medium&avatarId=10303&avatarType=issuetype'
    }
  };
  return function (input) {
    switch (input.toUpperCase()) {
      case type.STORY.name:
        return type.STORY.url;
      case type.BUG.name:
        return type.BUG.url;
    }
  }
});