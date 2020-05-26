import {Component, OnInit} from '@angular/core';
import {AuthService} from "../../services/auth/auth.service";
import {FormBuilder, Validators} from "@angular/forms";
import {Router} from "@angular/router";


@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  loginForm = this.formBuilder.group({
    email: this.formBuilder.control('', [Validators.required, Validators.email]),
    password: this.formBuilder.control('', [Validators.required]),
  });


  constructor(
    private  authService: AuthService,
    private formBuilder: FormBuilder,
    private router: Router,
  ) {
  }

  ngOnInit() {
  }

  login(customerData) {
    const usr = this.authService.login(customerData.email, customerData.password).
    then(err => {
      if (err.code === 'auth/wrong-password'){
        this.loginForm.setErrors({invalidCredentials: true})
      }
    });
  }

  loginInWithGoogle(){
    this.authService.loginWithGoogle();
  }

}
