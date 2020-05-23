import {Component, OnInit} from '@angular/core';
import {AuthService} from "../../services/auth/auth.service";
import {FormBuilder} from "@angular/forms";
import {error} from "util";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  checkoutForm;


  constructor(
    private  authService: AuthService,
    private formBuilder: FormBuilder,
  ) {
    this.checkoutForm = this.formBuilder.group({
      email: '',
      password: ''
    });
  }

  ngOnInit() {
  }

  login(customerData) {
    this.checkoutForm.reset();
    const usr = this.authService.login(customerData.email, customerData.password).
    then(value => console.log(value)).
    catch(err => console.log('Error ', err));
  }

  loginInWithGoogle(){
    this.authService.loginWithGoogle()
  }

}
