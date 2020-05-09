import {Component, OnInit} from '@angular/core';
import {UserModel, FriendRequestModel, FriendsModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {UserViewModel} from "./friendRequestView";
import {friendsColumnsDefinition} from "./table-definitions";

@Component({
  selector: 'app-my-profile',
  templateUrl: './my-profile.component.html',
  styleUrls: ['./my-profile.component.css']
})
export class MyProfileComponent implements OnInit {
  me: UserModel;
  IsendFriendRequests: UserViewModel[];
  IreceivedFriendRequests: UserViewModel[];
  myFriends: FriendsModel[];
  myFriendsView: UserViewModel[];
  newFriendEmail: string;

  friendsColumns = friendsColumnsDefinition;

  constructor(
    private userService: UserService,
  ) {
    this.IsendFriendRequests = [];
    this.IreceivedFriendRequests = [];
    this.myFriends = [];
    this.myFriendsView = [];
  }

  ngOnInit() {
    this.userService.getMe().subscribe(
      user => {
        this.me = user;
        this.userService.getMyFriendRequests().subscribe(
          requests => {
            if (requests != null) {
              this.fetchUserDetailsForFriendRequests(requests)
            }
          },
          err => {
            console.log('Error fetching my requests ', err)
          },
        )
      })

    this.userService.getMyFriends().subscribe(
      friends => {
        if (friends != null) {
          this.myFriends = friends;
          this.fetchUserDetailsForFriends(friends)
        }
      },
      error => console.log('Error GetMyFriends ', error)
    )
  }

  removeFriend(userId: string) {
    this.userService.removeFriend(userId).subscribe(
      () => {
        const index = this.myFriendsView.findIndex(u => u.UserId === userId)
        this.myFriendsView.splice(index, 1)
      }
    )

  }

  sendFriendRequest() {
    if (this.newFriendEmail != '') {
      this.userService.requestFriend(this.newFriendEmail).subscribe(
        newRequest => {
          this.userService.getUserByEmail(this.newFriendEmail).subscribe(
            user => {
              this.pushUserToFriendRequests([newRequest], [user])
              this.newFriendEmail = ''
            },
            error => console.log('Err getUser by email'))
        },)
    }
  }

  acceptFriendRequest(userId: string) {
    this.userService.acceptFriend(userId).subscribe(
      () => {
        const index = this.IreceivedFriendRequests.findIndex(u => u.UserId === userId)
        this.myFriendsView.push(this.IreceivedFriendRequests[index])
        this.IreceivedFriendRequests.splice(index, 1)
      },
      error => console.log('An error AcceptFreindRequest ', error)
    )
  }

  denyFriendRequest(userId: string) {
    this.userService.denyFriend(userId).subscribe(
      () => {
        const index = this.IreceivedFriendRequests.findIndex(u => u.UserId === userId)
        this.IreceivedFriendRequests.splice(index, 1)
      },
      error => console.log('An error denyFreindRequest ', error)
    )
  }

  cancelFriendRequest(userId: string) {
    this.userService.denyFriend(userId).subscribe(
      () => {
        const index = this.IsendFriendRequests.findIndex(u => u.UserId === userId)
        this.IsendFriendRequests.splice(index, 1)
      },
      error => console.log('An error denyFreindRequest ', error)
    )
  }


  private fetchUserDetailsForFriendRequests(requests: FriendRequestModel[]) {
    const userIds = requests.map(r => {
      if (r.User1Id !== this.me.Id)
        return r.User1Id;
      if (r.User2Id !== this.me.Id)
        return r.User2Id;
    });
    if (userIds.length !== 0) {
      this.userService.getUsersDetails(userIds).subscribe(
        users => {
          console.log(users)
          this.pushUserToFriendRequests(requests, users)
        },
        err => console.log('Error fetching user details ', err)
      );
    }
  }

  private fetchUserDetailsForFriends(friends: FriendsModel[]) {
    const userIds = friends.map(friend => {
      if (friend.User1Id !== this.me.Id)
        return friend.User1Id;
      if (friend.User2Id !== this.me.Id)
        return friend.User2Id;
    });
    if (userIds.length !== 0) {
      this.userService.getUsersDetails(userIds).subscribe(
        users => {
          this.pushUserToMyFriends(users)
        },
        err => console.log('Error fetching user details ', err)
      );
    }
  }

  private pushUserToMyFriends(users: UserModel[]) {
    this.myFriends.map(friendship => {
      users.map(user => {
          if (user.Id === friendship.User1Id || user.Id === friendship.User2Id) {
            this.myFriendsView.push(
              {
                UserId: user.Id,
                FirstName: user.FirstName,
                LastName: user.LastName,
                Email: user.Email,
                Avatar: user.Avatar,
                CreatedAt: friendship.CreatedAt,
              }
            )
          }
        }
      )
    })
  }

  private pushUserToFriendRequests(requests: FriendRequestModel[],users: UserModel[]) {
    requests.map(request => {
      users.map(user => {
        if (request.User1Id == user.Id) {
          this.IreceivedFriendRequests.push(
            {
              UserId: user.Id,
              FirstName: user.FirstName,
              LastName: user.LastName,
              Email: user.Email,
              Avatar: user.Avatar,
              CreatedAt: request.CreatedAt,
            }
          )
        } else if (request.User2Id == user.Id) {
          this.IsendFriendRequests.push(
            {
              UserId: user.Id,
              FirstName: user.FirstName,
              LastName: user.LastName,
              Email: user.Email,
              Avatar: user.Avatar,
              CreatedAt: request.CreatedAt,
            }
          )
        }
      })
    })
  }

}
