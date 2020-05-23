import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {ApisModel} from "../../api/api-urls";
import {LoadImageProps, SetLoadedImageProps,} from 'ng-lazyload-image';

@Component({
  selector: 'app-image',
  template: `
    <img id="test" #element [defaultImage]="defaultImage" [lazyLoad]="url">
  `,
  styleUrls: ['./lazy-load-img.component.css']
})
export class LazyLoadImgComponent implements OnInit {
  @ViewChild('element') el: HTMLImageElement;

  url: string;

  @Input() image: string;
  @Input() service: string;

  exampleBinUrl: string = 'https://static.independent.co.uk/s3fs-public/thumbnails/image/2019/04/15/17/alto-habitat-bin.jpg?w660';
  defaultImage: string = 'https://cdn.onlinewebfonts.com/svg/img_258083.png';

  constructor(
  ) {
  }

  ngOnInit() {
    if (this.service === ApisModel.user) {
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.user + '/load/' + this.image
      this.checkEmpty(this.defaultImage)
    }
    else if (this.service === ApisModel.society){
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.society + '/load/' + this.image
      this.checkEmpty(this.defaultImage)
    }
    else if (this.service === ApisModel.collection){
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.collection + '/load/' + this.image
      this.defaultImage = this.exampleBinUrl
      this.checkEmpty(this.defaultImage)
    }
    else if (this.service === ApisModel.trash){
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.trash + '/load/' + this.image
      this.defaultImage = this.exampleBinUrl
      this.checkEmpty(this.defaultImage)
    }

    console.log('URL: ', this.url)
  }

  checkEmpty(exampleUrl: string){
    if (this.image === ''){
      this.url = exampleUrl
    }
  }

  //setting downloaded image does not work

  // async ngAfterViewInit() {
  //   if (this.image === '' || this.image == null){
  //     this.el.src = this.defaultImage
  //     return
  //   }
  //
  //   if (this.el) {
  //     const imageBlob = await this.loadImage({
  //       element: this.el,
  //       imagePath: this.url
  //     })
  //
  //     console.log('blob is: ', imageBlob)
  //     this.el.src = imageBlob
  //     this.setLoadedImage({element: this.el, imagePath: imageBlob})
  //   }
  // }
  //
  // async loadImage({imagePath}: LoadImageProps) {
  //   const token = this.authService.getToken()
  //   return await fetch(imagePath, {
  //     headers: {
  //       Authorization: 'Bearer ' + token
  //     }
  //   }).then(res => res.blob()).then(blob => URL.createObjectURL(blob))
  // }
  //
  // setLoadedImage({element, imagePath, useSrcset}: SetLoadedImageProps) {
  //   console.log(this.el.currentSrc)
  //   this.el.src = imagePath
  //   console.log(this.el.currentSrc)
  // }

}
