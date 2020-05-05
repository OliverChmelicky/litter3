import { Component, OnInit } from '@angular/core';
import {NavigationItemModel} from "./navigationItem.model";
import {Observable} from "rxjs";
import {AuthService} from "../../services/auth/auth.service";

@Component({
  selector: 'app-nav-list',
  templateUrl: './nav-list.component.html',
  styleUrls: ['./nav-list.component.css']
})
export class NavListComponent implements OnInit {
  loggedInItems: NavigationItemModel[] = [
    {
      name: 'Report dump',
      url: '/report',
      icon:'delete',
    },
    {
      name: 'Societies',
      url: '/societies',
      icon:'delete',
    },
    {
      name: 'Map',
      url: '/map',
      icon:'my_location',
    },
    {
      name: 'Me',
      url: '/me',
      icon:'account_circle',
    },
  ];

  everyoneItems: NavigationItemModel[] = [
    {
      name: 'Report dump',
      url: '/report',
      icon:'delete',
    },
    {
      name: 'Societies',
      url: '/societies',
      icon:'delete',
    },
    {
      name: 'Map',
      url: '/map',
      icon:'my_location',
    },
    {
      name: 'Register',
      url: '/register',
      icon:'account_circle',
    },
    {
      name: 'Login',
      url: '/login',
      icon:'account_circle',
    },
  ]


  isLoggedIn$: Observable<boolean>;

  constructor(
    private authService: AuthService
  ) { }

  ngOnInit() {
    this.isLoggedIn$ = this.authService.isLoggedIn;
    this.authService.isLoggedIn.subscribe(isLogged => console.log(isLogged));
  }


  logout(){
    this.authService.logout()
  }

}
