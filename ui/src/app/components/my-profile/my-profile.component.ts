import {Component, OnInit} from '@angular/core';
import {UserModel, FriendRequestModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";

@Component({
  selector: 'app-my-profile',
  templateUrl: './my-profile.component.html',
  styleUrls: ['./my-profile.component.css']
})
export class MyProfileComponent implements OnInit {
  me: UserModel;
  friendRequests: FriendRequestModel[];
  //In future I can split into mine and other requests
  // sendFriendRequests: FriendRequestModel[];
  // ReceivedFriendRequests: FriendRequestModel[];
  newFriendEmail: string;

  constructor(
    private userService: UserService,
  ) {

  }

  ngOnInit() {
    this.userService.getMe().subscribe(
      user => {
        this.me = user;
        this.userService.getMyFriendRequests().subscribe(
          requests => {
            this.friendRequests = requests;
            const userIds = requests.map(r => {
              if (r.User1Id !== this.me.Id)
                return r.User1Id;
              if (r.User2Id !== this.me.Id)
                return r.User2Id;
            });
            this.userService.getUsersDetails(userIds).subscribe(
              users => console.log(users)
            );
          },
          err => {
            console.log('Error fetching my requests ',err)
          },

          //opt2 nablbo na debila getuj jednotlive ids
        )
    })
  }

  sendFriendRequest() {
    if (this.newFriendEmail != '') {
      this.userService.requestFriend(this.newFriendEmail).subscribe(
        reqest => console.log(reqest),
        err => console.log(err)
      )
    }
  }

}
