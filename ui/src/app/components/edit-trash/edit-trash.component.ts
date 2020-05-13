import { Component, OnInit } from '@angular/core';
import {TrashModel, TrashTypeBooleanValues} from "../../models/trash.model";
import {accessibilityChoces} from "../../models/accessibilityChocies";
import {ActivatedRoute, Router} from "@angular/router";
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
  trashTypeBool: TrashTypeBooleanValues;

  trashForm = this.formBuilder.group({
    lat: [''],
    lng: [''],
    size: 0,
    cleaned: false,

    trashTypeHousehold: false,
    trashTypeAutomotive: false,
    trashTypeConstruction: false,
    trashTypePlastics: false,
    trashTypeElectronic: false,
    trashTypeGlass: false,
    trashTypeMetal: false,
    trashTypeDangerous: false,
    trashTypeCarcass: false,
    trashTypeOrganic: false,
    trashTypeOther: false,

    accessibility: [''],
    description: '',
    anonymously: [false],
  });

  accessibilityChoices = accessibilityChoces;

  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private router: Router,
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
          this.printSize()
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

  onSave() {
    this.trash = {
      Id: this.trash.Id,
      Cleaned: this.trashForm.value["cleaned"],
      Size: this.sizeView,
      Accessibility: this.trashForm.value["accessibility"],
      TrashType: this.changeTrashTypeToInt(),
      Location: this.trash.Location,
      Description: this.trashForm.value["description"],
      FinderId: this.trash.FinderId,
      CreatedAt: this.trash.CreatedAt,
    }

    this.trashService.updateTrash(this.trash).subscribe(
      () => this.location.back()
    )
  }

  onDelete() {
    this.trashService.deleteTrash(this.trash.Id).subscribe(
      () => this.router.navigateByUrl('/map')
    )
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
    this.trashTypeBool = this.trashService.convertTrashTypeNumToBools(this.trash.TrashType);
  }

  printSize() {
    if (this.sizeValue == 0) {
      this.sizeView = 'unknown';
    }
    if (this.sizeValue == 1) {
      this.sizeView = 'bag';
    }
    if (this.sizeValue == 2) {
      this.sizeView = 'wheelbarrow';
    }
    if (this.sizeValue == 3) {
      this.sizeView = 'car';
    }
  }

  private changeTrashTypeToInt(): number {
    return this.trashService.convertTrashTypeBoolsToNums(
      {
        TrashTypeHousehold: !!this.trashForm.value.trashTypeHousehold,
        TrashTypeAutomotive: !!this.trashForm.value.trashTypeAutomotive,
        TrashTypeConstruction: !!this.trashForm.value.trashTypeConstruction,
        TrashTypePlastics: !!this.trashForm.value.trashTypePlastics,
        TrashTypeElectronic: !!this.trashForm.value.trashTypeElectronic,
        TrashTypeGlass: !!this.trashForm.value.trashTypeGlass,
        TrashTypeMetal: !!this.trashForm.value.trashTypeMetal,
        TrashTypeDangerous: !!this.trashForm.value.trashTypeDangerous,
        TrashTypeCarcass: !!this.trashForm.value.trashTypeCarcass,
        TrashTypeOrganic: !!this.trashForm.value.trashTypeOrganic,
        TrashTypeOther: !!this.trashForm.value.trashTypeOther,
      }
    );
  }
}
