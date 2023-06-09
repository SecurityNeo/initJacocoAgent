package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"initJacocoAgent/service"
	"io/ioutil"
	adminssionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func WebhookCallBack(c *gin.Context) {
	var body []byte
	if c.Request.Body != nil {
		if data, err := ioutil.ReadAll(c.Request.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		errMsg := "empty body"
		log.Warnf(errMsg)
		c.Status(http.StatusBadRequest)
		return
	}

	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" {
		errMsg := fmt.Sprintf("The Content-Type is %s,need application/json", contentType)
		log.Warnf(errMsg)
		c.Status(http.StatusUnsupportedMediaType)
		return
	}

	var admissionRsp *adminssionv1.AdmissionResponse
	// v1版本需要指定Kind和APIVersion，否则客户端会报错
	ar := adminssionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
	}
	if _, _, er := deserializer.Decode(body, nil, &ar); er != nil {
		errMsg := fmt.Sprintf("Can't decode body.error msg: %s", er.Error())
		log.Errorf(errMsg)
		admissionRsp = &adminssionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: er.Error(),
			},
		}
	} else {
		admissionRsp = Mutate(&ar)
	}
	admissionReview := adminssionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
	}
	if admissionRsp != nil {
		admissionReview.Response = admissionRsp
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}
	var conLog bytes.Buffer
	rsp, err := json.Marshal(admissionReview)
	if err != nil {
		conLog.WriteString(fmt.Sprintf("Can't encode response: %v", err))
		c.String(http.StatusInternalServerError, conLog.String())
		return
	}
	c.Writer.Write(rsp)
}

func Mutate(ar *adminssionv1.AdmissionReview) *adminssionv1.AdmissionResponse {
	req := ar.Request

	fmt.Printf("\n ----Begin fo admission for NS=[%v],Kind=[%v],Name=[%v]----", req.Namespace, req.Kind.Kind, req.Name)
	switch req.Kind.Kind {
	case "Deployment":
		var dp appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &dp); err != nil {
			errMsg := fmt.Sprintf("\nClould not unmarshal raw object: %v", err)
			fmt.Printf(errMsg)
			return &adminssionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		return service.MutateDeploy(&dp)
	case "StatefulSet":
		var sts appsv1.StatefulSet
		if err := json.Unmarshal(req.Object.Raw, &sts); err != nil {
			errMsg := fmt.Sprintf("\nClould not unmarshal raw object: %v", err)
			fmt.Printf(errMsg)
			return &adminssionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		return service.MutateSts(&sts)
	default:
		msg := fmt.Sprintf("\n Do not support for this kind of resource %v", req.Kind.Kind)
		fmt.Printf(msg)
		return &adminssionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}
}
