import { Component, OnInit } from '@angular/core';
import {membersColumnsSocietyEditDefinition, roles} from "../society-details/table-definitions";
import {ActivatedRoute, Router} from "@angular/router";
import {SocietyService} from "../../services/society/society.service";
import {UserService} from "../../services/user/user.service";
import {MemberModel, SocietyModel} from "../../models/society.model";
import {UserInSocietyModel, UserModel} from "../../models/user.model";
import {FormBuilder} from "@angular/forms";
import {FileuploadService} from "../../services/fileupload/fileupload.service";
import {MatSelectChange} from "@angular/material/select";

@Component({
  selector: 'app-edit-society',
  templateUrl: './edit-society.component.html',
  styleUrls: ['./edit-society.component.css']
})
export class EditSocietyComponent implements OnInit {
  membersColumnsDef = membersColumnsSocietyEditDefinition;
  changeMemberPermission: MemberModel[] = []
  members: UserInSocietyModel[];
  origMembers: UserInSocietyModel[];

  roles = roles;
  society: SocietyModel;
  fd: FormData = new FormData();

  adminsMembers: MemberModel[];
  isAdmin: boolean = false;
  me: UserModel;

  societyForm = this.formBuilder.group({
    description: '',
    name: '',
  });


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private societyService: SocietyService,
    private userService: UserService,
    private formBuilder: FormBuilder,
    private fileuploadService: FileuploadService,
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.societyService.getSociety(params.get('societyId')).subscribe(
        society => {
            this.society = society
            console.log(society.Avatar)
          this.userService.getMe().subscribe(
            res => {
              this.me = res
              this.societyService.getSocietyMembers(this.society.Id).subscribe(m => {
                this.getMembers(m)
                this.adminsMembers = m.filter( mem => mem.Permission === 'admin')
                this.adminsMembers.map( a => {
                    if (a.UserId === this.me.Id)
                      this.isAdmin = true;
                })
              })
            })
        })
    });
  }

  private getMembers(membersPermissions: MemberModel[]) {
    const membIds = membersPermissions.map( m => m.UserId)
    this.userService.getUsersDetails(membIds).subscribe( m => {
      this.members = [];
      m.map( usr =>
        {
          const roleOfUser = membersPermissions.map( membership =>
            {if (membership.UserId === usr.Id) {
              return membership
            }
          })
          if (roleOfUser.length === 1) {
            this.members.push({
              user: usr,
              role: roleOfUser[0].Permission,
            })
          }

        })
      this.origMembers = []
      for (let i = 0; i < this.members.length; i++) {
        this.origMembers.push(
          {
            user: this.members[i].user,
            role: this.members[i].role,
          }
        )
      }
    })
  }

  memberPermissionChange(event: MatSelectChange, i: number) {
    if (event.value === this.origMembers[i].role) {
      const index = this.changeMemberPermission.findIndex(u => u.UserId === this.members[i].user.Id)
      this.changeMemberPermission.splice(index, 1)
      console.log('reverted same')
    } else {
      const exists = this.changeMemberPermission.filter( mem => mem.UserId === this.members[i].user.Id)
      if (exists.length !== 0) {
        const index = this.changeMemberPermission.findIndex(u => u.UserId === this.members[i].user.Id)
        this.changeMemberPermission.splice(index, 1)
        console.log('reverted ine')
      }

      this.changeMemberPermission.push({
        UserId: this.members[i].user.Id,
        SocietyId: this.society.Id,
        Permission: event.value.toString(),
        CreatedAt: new Date(),  //server does not use this property
      })
      console.log('zmena vykonana')
    }
  }

  onMemberPermissionAcceptChanges() {
    if (this.changeMemberPermission.length > 0) {
      this.societyService.changePermissions(this.changeMemberPermission).subscribe()
    }
  }

  onUpdate() {
    if (this.societyForm.value['name'] !== '') {
      this.society.Name = this.societyForm.value['name']
    }
    if (this.societyForm.value['description'] !== '' ) {
      this.society.Description = this.societyForm.value['description']
    }
    console.log('idem volat society')
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
    this.societyService.removeUser(id, this.society.Id).subscribe(
      () => {
        if (this.me.Id === id) {
          this.router.navigateByUrl('map')
        } else {
          const index = this.members.findIndex( m => m.user.Id === id)
          this.members.splice(index, 1)
        }
      }
    )
  }

  onDelete() {
    //testni delete society
  }
}
