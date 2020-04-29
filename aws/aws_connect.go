package aws_test
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "fmt"
    "github.com/aws/aws-sdk-go/service/efs"
)

// Function to Create AWS EBS Volume
func CreateEBSVolume(ec2Instance *ec2.EC2, size int64, az, volType string ) {
    input := &ec2.CreateVolumeInput{
        AvailabilityZone: aws.String(az),
        Size:             aws.Int64(size),
        VolumeType:       aws.String(volType),
    }
    result_vol, err_vol := ec2Instance.CreateVolume(input)
    if err_vol != nil {
        if aerr, ok := err_vol.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            fmt.Println(err_vol.Error())
        }
    }
    fmt.Println(result_vol)
}


// Function to Delete AWS EBS Volume
func DeleteEBSVolume(ec2Instance *ec2.EC2, vol_id string) {
    deleteVolInput := &ec2.DeleteVolumeInput{
	VolumeId: aws.String(vol_id),
	}
    result_vol, err_vol := ec2Instance.DeleteVolume(deleteVolInput)
    if err_vol != nil {
        if aerr, ok := err_vol.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            fmt.Println(err_vol.Error())
        }
    }
    fmt.Println(result_vol)
}


// Function to get the details of all instances of AWS EC2
func DescribeInstance(ec2Instance *ec2.EC2) {
    // Call to get detailed information on each instance
    result, err := ec2Instance.DescribeInstances(nil)
    if err != nil {
        fmt.Println("Error: \n %v", err)
    } else {
        fmt.Println("Result:\n %v", result)
    }

}

// Function to stop or start Instances of EC2
func Start_stop_instances(ec2Instance *ec2.EC2, instance_id string, op_name string) {
    switch op_name {
        case "stop":
             input := &ec2.StopInstancesInput {
                          InstanceIds : []*string{
                                             aws.String(instance_id),
                                         },
                      }
             result, err := ec2Instance.StopInstances(input)
             if err != nil {
                 fmt.Println("Error in stopping the instance", err)
             } else {
		 fmt.Println("Success:", result.StoppingInstances)
             }
        case "start":
             input := &ec2.StartInstancesInput {
                          InstanceIds : []*string{
                                             aws.String(instance_id),
                                         },
                      }
             result, err := ec2Instance.StartInstances(input)
             if err != nil {
                 fmt.Println("Error in starting the instance", err)
             } else {
		 fmt.Println("Success:", result.StartingInstances)
             }
        default:
             fmt.Println("No input provided")
    }
}


// Function to attach a EBS volume to the EC2 instance
func AttachVol(ec2Instance *ec2.EC2, instance_id, vol_id string) {
    Start_stop_instances(ec2Instance, instance_id, "stop")
    input := &ec2.AttachVolumeInput{
         Device: aws.String("/dev/sdf"),
         InstanceId: aws.String(instance_id),
         VolumeId: aws.String(vol_id),
    }
    result, err := ec2Instance.AttachVolume(input)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println(result)
    }
    Start_stop_instances(ec2Instance, instance_id, "start")
}

// Create session with IAM AK-SK and the region
func CreateSession(region, ak, sk string) (*ec2.EC2) {
    sess,err := session.NewSession(&aws.Config{
            Region: aws.String(region),
            Credentials: credentials.NewStaticCredentials(ak, sk, ""),
    })
    if err != nil {
	fmt.Println("Error creating the session")
	return nil
    }
    ec2Instance := ec2.New(sess)
    return ec2Instance
}

func CreateEFSSession(region, ak, sk string) (*efs.EFS) {
    mySession := session.Must(session.NewSession(&aws.Config{
            Region: aws.String(region),
            Credentials: credentials.NewStaticCredentials(ak, sk, ""),
    }))

    // Create a EFS client from just a session.
    svc := efs.New(mySession)
    return svc
}

// Create a FS with the given name
func CreateEFS(efsInstance *efs.EFS, fsname string) {
    input := &efs.CreateFileSystemInput{
             CreationToken:   aws.String("tokenstring"),
             PerformanceMode: aws.String("generalPurpose"),
             Tags: []*efs.Tag{
                 {
                 Key:   aws.String("Name"),
                 Value: aws.String(fsname),
                   },
                 },
             }
    result, err := efsInstance.CreateFileSystem(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
             switch aerr.Code() {
             case efs.ErrCodeBadRequest:
                fmt.Println(efs.ErrCodeBadRequest, aerr.Error())
             case efs.ErrCodeInternalServerError:
                fmt.Println(efs.ErrCodeInternalServerError, aerr.Error())
             case efs.ErrCodeFileSystemAlreadyExists:
                fmt.Println(efs.ErrCodeFileSystemAlreadyExists, aerr.Error())
             case efs.ErrCodeFileSystemLimitExceeded:
                fmt.Println(efs.ErrCodeFileSystemLimitExceeded, aerr.Error())
             case efs.ErrCodeInsufficientThroughputCapacity:
                fmt.Println(efs.ErrCodeInsufficientThroughputCapacity, aerr.Error())
             case efs.ErrCodeThroughputLimitExceeded:
                fmt.Println(efs.ErrCodeThroughputLimitExceeded, aerr.Error())
             default:
                fmt.Println(aerr.Error())
            }
         } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return
    }

    fmt.Println(result)

}

// Delete a EFS with the ID
func DeleteEFS(efsInstance *efs.EFS, fsid string) {
    input := &efs.DeleteFileSystemInput{
                 FileSystemId: aws.String(fsid),
             }

    result, err := efsInstance.DeleteFileSystem(input)
    if err != nil {
	    if aerr, ok := err.(awserr.Error); ok {
		    switch aerr.Code() {
			    case efs.ErrCodeBadRequest:
				    fmt.Println(efs.ErrCodeBadRequest, aerr.Error())
			    case efs.ErrCodeInternalServerError:
				    fmt.Println(efs.ErrCodeInternalServerError, aerr.Error())
		            case efs.ErrCodeFileSystemNotFound:
				    fmt.Println(efs.ErrCodeFileSystemNotFound, aerr.Error())
			    case efs.ErrCodeFileSystemInUse:
			            fmt.Println(efs.ErrCodeFileSystemInUse, aerr.Error())
			    default:
				    fmt.Println(aerr.Error())
			    }
		    } else {
			    // Print the error, cast err to awserr.Error to get the Code and
			    // Message from an error.
			    fmt.Println(err.Error())
		    }
		    return
    }

    fmt.Println(result)
}
