import {Component, OnInit} from '@angular/core';
import {AuthService} from "../../services/auth/auth.service";
import {FormBuilder} from "@angular/forms";

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {
  errorMessage: string;
  successMessage: string;
  registerForm;

  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
  ) {
  }

  ngOnInit() {
    this.registerForm = this.formBuilder.group({
      email: '',
      password: ''
    });
  }

  tryRegiser(value) {
    this.authService.register(value)
      .then(() => {
        this.errorMessage = null;
        this.successMessage = "Your account has been created";
      }, err => {
        this.errorMessage = err.message;
        this.successMessage = null;
      })
  }


}
