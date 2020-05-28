import {Component, OnInit} from '@angular/core';
import {AuthService} from "../../services/auth/auth.service";
import {FormBuilder, Validators} from "@angular/forms";
import {Router} from "@angular/router";

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {
  errorMessage: string;
  registerForm;
  userExists: boolean;

  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.registerForm = this.formBuilder.group({
      email: this.formBuilder.control('', [Validators.required, Validators.email]),
      password: this.formBuilder.control('', [Validators.required, Validators.minLength(6)]),
    });
  }

  tryRegister(value) {
    this.authService.register(value)
      .then(() => {
        this.errorMessage = null;
        this.router.navigateByUrl('map');
      }, err => {
        this.errorMessage = err.message;
        this.userExists = true;
      })
  }

  registerWithGoogle() {
    this.authService.loginWithGoogle().then(() => {
      this.router.navigateByUrl('map');
    })
  }


}
