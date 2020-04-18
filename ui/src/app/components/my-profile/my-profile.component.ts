import {Component, OnInit} from '@angular/core';
import {UserModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";

@Component({
  selector: 'app-my-profile',
  templateUrl: './my-profile.component.html',
  styleUrls: ['./my-profile.component.css']
})
export class MyProfileComponent implements OnInit {
  me: UserModel;

  constructor(
    private userService: UserService,
  ) {

  }

  ngOnInit() {
    this.userService.getMe().subscribe(user => this.me = user)
  }

}
