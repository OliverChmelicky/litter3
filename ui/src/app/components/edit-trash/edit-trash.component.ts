import { Component, OnInit } from '@angular/core';
import {TrashModel} from "../../models/trash.model";
import {accessibilityChoces} from "../../models/accessibilityChocies";
import {ActivatedRoute} from "@angular/router";
import {Location} from '@angular/common';
import {TrashService} from "../../services/trash/trash.service";
import {FormBuilder} from "@angular/forms";

@Component({
  selector: 'app-edit-trash',
  templateUrl: './edit-trash.component.html',
  styleUrls: ['./edit-trash.component.css']
})
export class EditTrashComponent implements OnInit {
  trashId: string;
  trash: TrashModel;
  sizeView: string;
  sizeValue: number;
  fd: FormData = new FormData();

  trashForm = this.formBuilder.group({
    lat: [''],
    lng: [''],
    size: [1],

    trashTypeHousehold: [''],
    trashTypeAutomotive: [''],
    trashTypeConstruction: [''],
    trashTypePlastics: [''],
    trashTypeElectronic: [''],
    trashTypeGlass: [''],
    trashTypeMetal: [''],
    trashTypeDangerous: [''],
    trashTypeCarcass: [''],
    trashTypeOrganic: [''],
    trashTypeOther: [''],

    accessibility: [''],
    description: [''],
    anonymously: [''],
  });

  accessibilityChoices = accessibilityChoces;

  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private trashService: TrashService,
    private formBuilder: FormBuilder,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.trashId = params.get('id');
      this.trashService.getTrashById(this.trashId).subscribe(
        trash => {
          this.trash = trash
          this.convertSizeToNumber(trash.Size)
          this.convertTrashTypeToBools()
        }
      )
    });
  }

  convertSizeToNumber(size: string) {
    if (size === 'unknown') {
      this.sizeValue = 0;
    }
    if (size == 'bag') {
      this.sizeValue = 1;
    }
    if (size == 'wheelbarrow') {
      this.sizeValue = 2;
    }
    if (size == 'car') {
      this.sizeValue = 3;
    }
  }

  onSubmit() {
    this.trashService.updateTrash(this.trash).subscribe()
  }

  onDelete() {
    this.trashService.deleteTrash(this.trash.Id).subscribe()
  }

  onGoBack() {
    this.location.back()
  }

  onFileSelected(event) {
    this.fd.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.fd.append("files", event.target.files[i], event.target.files[i].name);
    }
  }

  private convertTrashTypeToBools() {
    //TODO
    console.log('pohraj sa s bitikmi')
  }
}
