import {Component} from '@angular/core';
import {AuthService} from "./services/auth/auth.service";
import {Observable} from "rxjs";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'ui';

  isLoggedIn$: Observable<boolean>;

  constructor(
    private authService: AuthService
  ) {
    this.isLoggedIn$ = this.authService.isLoggedIn;
    this.authService.isLoggedIn.subscribe();
  }

  logout() {
    this.authService.logout()
  }
}
