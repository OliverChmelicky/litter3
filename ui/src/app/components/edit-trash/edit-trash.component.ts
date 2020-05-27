import {Component, Inject, OnInit} from '@angular/core';
import {
  defaultTrashModel,
  defaultTrashTypeBooleanValues,
  TrashModel,
  TrashTypeBooleanValues
} from "../../models/trash.model";
import {accessibilityChoces} from "../../models/accessibilityChocies";
import {ActivatedRoute, Router} from "@angular/router";
import {Location} from '@angular/common';
import {TrashService} from "../../services/trash/trash.service";
import {FormBuilder} from "@angular/forms";
import {FileuploadService} from "../../services/fileupload/fileupload.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {MarkerModel} from "../google-map/Marker.model";
import {GoogleMap} from "@agm/core/services/google-maps-types";

export interface DialogData {
  Url: string;
}

@Component({
  selector: 'app-edit-trash',
  templateUrl: './edit-trash.component.html',
  styleUrls: ['./edit-trash.component.css']
})
export class EditTrashComponent implements OnInit {
  trash: TrashModel = defaultTrashModel;
  sizeView: string;
  sizeValue: number;
  fd: FormData = new FormData();
  trashTypeBool: TrashTypeBooleanValues = defaultTrashTypeBooleanValues;

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
  });
  accessibilityChoices = accessibilityChoces;
  errorMessage: string;

  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private router: Router,
    private trashService: TrashService,
    private formBuilder: FormBuilder,
    private fileuploadService: FileuploadService,
    public openImageDialog: MatDialog,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const trashId = params.get('id');
      this.trashService.getTrashById(trashId).subscribe(
        trash => {
          this.trash = trash
          this.trashForm.controls['cleaned'].setValue(trash.Cleaned) //needed for correct using of checkbox
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
      () => {
        if (this.fd.getAll('files').length !== 0) {
          this.fileuploadService.uploadTrashImages(this.fd, this.trash.Id).subscribe(
            () => {
              this.fd.delete('files')
              this.location.back()
            })
        } else {
          this.location.back()
        }
      }
    )
  }

  onDelete() {
    this.trashService.deleteTrash(this.trash.Id).subscribe(
      () => this.router.navigateByUrl('/map'),
      error => {
        this.errorMessage = 'You cannot delete trash which has collections or events!'
      }
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

  onDeleteImage(imageUrl: string) {
    this.trashService.deleteTrashImage(imageUrl, this.trash.Id).subscribe()
    const index = this.trash.Images.findIndex(i => i.Url === imageUrl)
    this.trash.Images.splice(index, 1)
  }

  onOpenImage(url: string): void {
    const dialogRef = this.openImageDialog.open(ImageDialogComponent, {
      width: '800px',
      data: {
        Url: url,
      }
    });
  }
}

@Component({
  selector: 'app-image-dialog',
  templateUrl: '../image-dialog/image-dialog.component.html',
  styleUrls: ['../image-dialog/image-dialog.component.css']
})
export class ImageDialogComponent {

  constructor(public dialogRef: MatDialogRef<ImageDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: DialogData) {
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

}
