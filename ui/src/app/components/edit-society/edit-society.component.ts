import {Component, Inject, OnInit} from '@angular/core';
import {membersColumnsSocietyEditDefinition, roles} from "../society-details/table-definitions";
import {ActivatedRoute, Router} from "@angular/router";
import {SocietyService} from "../../services/society/society.service";
import {UserService} from "../../services/user/user.service";
import {DefaultSociety, MemberModel, SocietyModel} from "../../models/society.model";
import {UserInSocietyModel, UserModel} from "../../models/user.model";
import {FormBuilder} from "@angular/forms";
import {FileuploadService} from "../../services/fileupload/fileupload.service";
import {MatSelectChange} from "@angular/material/select";
import {ApisModel} from "../../api/api-urls";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {animate, state, style, transition, trigger} from "@angular/animations";
import {MatTableDataSource} from "@angular/material/table";
import {MarkerModel} from "../google-map/Marker.model";

export interface DialogData {
  header: string,
  commandName: string,
}

@Component({
  selector: 'app-edit-society',
  templateUrl: './edit-society.component.html',
  styleUrls: ['./edit-society.component.css'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({height: '0px', minHeight: '0'})),
      state('expanded', style({height: '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class EditSocietyComponent implements OnInit {
  membersColumnsDef = membersColumnsSocietyEditDefinition;
  changeMemberPermission: MemberModel[] = []
  members: UserInSocietyModel[] = [];
  origMembers: UserInSocietyModel[] = [];

  roles = roles;
  society: SocietyModel = DefaultSociety;
  fd: FormData = new FormData();

  adminsMembers: MemberModel[];
  isAdmin: boolean = false;
  me: UserModel;

  societyForm = this.formBuilder.group({
    description: '',
    name: '',
  });
  assignSociety: string = ApisModel.society;
  assignUser: string = ApisModel.user;


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private societyService: SocietyService,
    private userService: UserService,
    private formBuilder: FormBuilder,
    private fileuploadService: FileuploadService,
    public removeUserDialog: MatDialog,
  ) { }

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      this.societyService.getSociety(params.get('societyId')).subscribe(
        society => {
          this.society = society

          let adminNumber = 0
            this.society.MemberRights.map( m => {
              if (m.Permission === 'admin') {
                adminNumber += 1;
              }
            })

          society.Users.map( user => {
            this.society.MemberRights.map( right => {
              if (user.Id === right.UserId) {
                let showRemove = true
                if (adminNumber === 1 && right.Permission === 'admin') {
                  showRemove = false
                }

                this.members.push({
                  user: user,
                  role: right.Permission,
                  showRemove: showRemove
                })
              }
            })
          })
          this.members.map( m => this.origMembers.push({
            user: m.user,
            role: m.role,
            showRemove: m.showRemove
          }))
          this.refreshMembersTable()

          this.userService.getMe().subscribe(
            res => {
              this.me = res
              this.adminsMembers = this.society.MemberRights.filter( mem => mem.Permission === 'admin')
              this.adminsMembers.map( a => {
                if (a.UserId === this.me.Id)
                  this.isAdmin = true;
              })
            })
        })
    });
  }

  memberPermissionChange(event: MatSelectChange, i: number) {
    console.log(this.origMembers[i])
    if (event.value === this.origMembers[i].role) {
      //back to the same permission
      const index = this.changeMemberPermission.findIndex(u => u.UserId === this.members[i].user.Id)
      this.changeMemberPermission.splice(index, 1)
      console.log('Changed permission: ',this.changeMemberPermission)
      return
    } else {
      //find old permission
      const exists = this.changeMemberPermission.filter( mem => mem.UserId === this.members[i].user.Id)
      console.log('exists? ', exists)
      if (exists.length !== 0) {
        //remove
        const index = this.changeMemberPermission.findIndex(u => u.UserId === this.members[i].user.Id)
        this.changeMemberPermission.splice(index, 1)
      }

      //add new
      this.changeMemberPermission.push({
        UserId: this.members[i].user.Id,
        SocietyId: this.society.Id,
        Permission: event.value.toString(),
        CreatedAt: new Date(),  //server does not use this property
      })
      console.log('Changed permission: ',this.changeMemberPermission)
    }
  }

  onMemberPermissionAcceptChanges() {
    console.log(this.changeMemberPermission)
    if (this.changeMemberPermission.length > 0) {
      this.societyService.changePermissions(this.changeMemberPermission).subscribe(
        () => {
          this.changeMemberPermission = [];

          this.origMembers = []
          this.members.map( m => this.origMembers.push({
            user: m.user,
            role: m.role,
            showRemove: m.showRemove
          }))
        },
        error => console.log(error)
      )
    }
  }

  onUpdate() {
    console.log(this.societyForm.value)
    if (this.societyForm.value['name'] !== '') {
      this.society.Name = this.societyForm.value['name']
    }
    if (this.societyForm.value['description'] !== '' ) {
      this.society.Description = this.societyForm.value['description']
    }
    this.societyService.updateSociety(this.society).subscribe()
    if (this.fd.getAll('file').length !== 0) {
      this.fileuploadService.uploadSocietyImage(this.fd, this.society.Id).subscribe(
        () => this.fd.delete('file')
      )
    }
  }

  onFileSelected(event) {
    this.fd.delete('file')
    this.fd.append("file", event.target.files[0], event.target.files[0].name);
  }

  removeUser(id: string) {
    const dialogRef = this.removeUserDialog.open(RemoveMemberComponent, {
      width: '800px',
      data: {
        header: 'Remove user?',
        commandName: 'REMOVE',
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log('Idem removnut: ',id)
        this.societyService.removeUser(id, this.society.Id).subscribe(
          () => {
            if (this.me.Id === id) {
              this.router.navigateByUrl('map')
            } else {
              const index = this.members.findIndex( m => m.user.Id === id)
              this.members.splice(index, 1)

              this.refreshMembersTable()
            }
          }
        )
      }
    });
  }

  onDelete() {
    const dialogRef = this.removeUserDialog.open(RemoveMemberComponent, {
      width: '800px',
      data: {
        header: 'Delete society?',
        commandName: 'DELETE',
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.societyService.deleteSociety(this.society.Id).subscribe(
          () => {
            this.router.navigateByUrl('map')
          }
        )
      }
    });
  }

  refreshMembersTable() {
    const newData = new MatTableDataSource<UserInSocietyModel>(this.members);
    this.members = []
    for (let i = 0; i < newData.data.length; i++) {
      this.members.push(newData.data[i])
    }
  }
}

@Component({
  selector: 'app-edit-soc-modal',
  templateUrl: './dialog/edit-soc-modal.component.html',
  //styleUrls: ['./dialog/edit-soc-modal.component.css]
})
export class RemoveMemberComponent {

  constructor(public dialogRef: MatDialogRef<RemoveMemberComponent>,
              @Inject(MAT_DIALOG_DATA) public data: DialogData) {
  }

  onNoClick(): void {
    this.dialogRef.close();
  }
}
