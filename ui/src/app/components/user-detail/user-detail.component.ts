import { Component, OnInit } from '@angular/core';
import {UserService} from "../../services/user/user.service";
import {UserModel} from "../../models/user.model";

@Component({
  selector: 'app-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.css']
})
export class UserDetailComponent implements OnInit {
user: UserModel;



  constructor(
    private userService: UserService,
  ) {

  }

  ngOnInit() {

  }


  testFetch(id :string) {
    const isRegistered = this.userService.getRegistered()
    console.log(isRegistered)
    if (isRegistered){
      this.user = isRegistered;
    } else {
      this.userService.getUser(id).
      subscribe(user => this.user = user);
    }

  }

}
